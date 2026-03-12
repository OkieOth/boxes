<%
    import re
    import yacg.model.model as model
    import yacg.templateHelper as templateHelper
    import yacg.model.modelFuncs as modelFuncs
    import yacg.util.stringUtils as stringUtils

    templateFile = 'golang_types.mako'
    templateVersion = '1.1.0'

    packageName = templateParameters.get('modelPackage','<<PLEASE SET modelPackage TEMPLATE PARAM>>')
    basePackageName = templateParameters.get('baseModelPackage','<<PLEASE SET baseModelPackage TEMPLATE PARAM>>')
    basePackageImport = templateParameters.get('baseModelPackageImport','')
    jsonTypesPackage = templateParameters.get('jsonTypesPackage','<<PLEASE SET jsonTypesPackage TEMPLATE PARAM>>')
    jsonSerialization = templateParameters.get('jsonSerialization',False)

    noExternalTypes = templateParameters.get('noExternalTypes',False)

    def printStarForJson(isJson):
        return "*" if isJson else ""

    def printAsteriskIfNotRequired(property):
        return  "" if property.required else "*"

    def printAsteriskIfRequired(property):
        return  "*" if property.required else ""

    def printAndIfNotRequired(property):
        return  "" if property.required else "&"

    def printAndIfRequired(property):
        return  "&" if property.required else ""

    def hasDictOrArrayAttribs(type):
        for property in type.properties:
            if property.isArray or isinstance(property.type, model.DictionaryType):
                return True
        return False

    def copyFunctionName(property):
        t = printGolangType(property.type, False, True, 0, False)
        i = t.rfind('.')
        if i != -1:
            return t[:i] + ".Copy" + t[i+1:]
        else:
            return "Copy"+t


    def printGolangType(typeObj, isArray, isRequired, arrayDimensions, forJson):
        ret = ''
        if typeObj is None:
            return '???'
        elif isinstance(typeObj, model.IntegerType):
            if typeObj.format == model.IntegerTypeFormatEnum.INT32:
                ret = 'int32'
            elif typeObj.format == model.IntegerTypeFormatEnum.UINT32:
                ret = 'uint32'
            else:
                ret = 'int'
        elif isinstance(typeObj, model.ObjectType):
            ret = 'interface{}'
        elif isinstance(typeObj, model.NumberType):
            if typeObj.format is None or typeObj.format == model.NumberTypeFormatEnum.DOUBLE:
                ret = 'float64'
            else:
                ret = 'float32'
        elif isinstance(typeObj, model.BooleanType):
            ret = 'bool'
        elif isinstance(typeObj, model.StringType):
            ret = 'string'
        elif isinstance(typeObj, model.BytesType):
            ret = 'byte'
        elif isinstance(typeObj, model.UuidType):
            ret = 'uuid.UUID'
        elif isinstance(typeObj, model.EnumType):
            ret = typeObj.name
        elif isinstance(typeObj, model.DateType):
            ret = 'time.Time'
        elif isinstance(typeObj, model.TimeType):
            ret = 'time.Time'
        elif isinstance(typeObj, model.DateTimeType):
            ret = 'time.Time'
        elif isinstance(typeObj, model.DictionaryType):
            ret = 'map[string]{}'.format(printGolangType(typeObj.valueType, False, True, 0, False))
        elif isinstance(typeObj, model.ComplexType):
            if (typeObj.domain == "base") and (not noExternalTypes):
                ret = "{}.{}".format(basePackageName, typeObj.name)
            else:
                ret = typeObj.name
        else:
            ret = '???'

        if isArray:
            ret = "[]{}".format(ret)

        if (not isRequired) and (not (property.isArray or isinstance(property.type, model.DictionaryType))):
            ret = "*{}".format(ret)

        return ret


    def printOmitemptyIfNeeded(property):
        if not property.required or property.isArray or isinstance(property.type, model.DictionaryType):
            return ",omitempty"
        else:
            return ""

    def printStarIfRequired(property):
        if property.required and (not property.isArray) and (not isinstance(property.type, model.DictionaryType)):
            return "*"
        else:
            return ""


    def getEnumDefaultValue(type):
        if type.default is not None:
            return secureEnumValues(type.default, type.name)
        else:
            return secureEnumValues(type.values[0], type.name)

    def secureEnumValues(value, typeName):
        valueName = stringUtils.toName(value)
        typeName = stringUtils.toName(typeName)
        return typeName + "_" + valueName
        #pattern = re.compile("^[0-9]")
        #return '_' + value if pattern.match(value) else value

    def isEnumDefaultValue(value, type):
        return getEnumDefaultValue(type) == secureEnumValues(value, type.name)

    def sanitizePropertyName(property):
        name = property.name
        if name == "type":
            return "type_"
        return name


    def needsBuilder(type):
        if hasattr(type, "properties"):
            for property in type.properties:
                if (property.type is not None) and (isinstance(property.type, model.DictionaryType)):
                    return True
                if property.isArray:
                    return True
                if isinstance(property.type, model.ComplexType):
                    b = needsBuilder(property.type)
                    if b:
                        return True
        return False


%>package ${packageName}

// Attention, this file is generated. Manual changes get lost with the next
// run of the code generation.
// created by yacg (template: ${templateFile} v${templateVersion})

import (
% if modelFuncs.isUuidContained(modelTypes):
    uuid "github.com/google/uuid"
% endif
% if modelFuncs.isAtLeastOneDateRelatedTypeContained(modelTypes):
    "time"
% endif
% if modelFuncs.hasEnumTypes(modelTypes):
    "encoding/json"
    "errors"
    "fmt"
% endif
% if basePackageImport != "":
    "${basePackageImport}"
% endif
)

% for type in modelTypes:
    % if modelFuncs.isEnumType(type):
        % if type.description != None:
/* ${templateHelper.addLineBreakToDescription(type.description,4)}
*/
        % endif
type ${type.name} string

const (
        % for value in type.values:
    ${secureEnumValues(value, type.name)} ${type.name} = "${value}"
        % endfor
)


func (s *${type.name}) UnmarshalJSON(data []byte) error {
    var value string
    if err := json.Unmarshal(data, &value); err != nil {
        return err
    }

    switch value {
        % for value in type.values:
    case "${value}":
        *s = ${secureEnumValues(value, type.name)}
        % endfor
    default:
		msg := fmt.Sprintf("invalid value for DDDDomainType: %s", value)
		return errors.New(msg)
    }
    return nil
}

    % endif

    % if hasattr(type, "properties"):
        % if type.description != None:
/* ${templateHelper.addLineBreakToDescription(type.description,4)}
*/
        % endif
type ${type.name} struct {
        % for property in type.properties:

            % if property.description != None:
    // ${property.description}
            % endif
    ${stringUtils.toUpperCamelCase(property.name)} ${printGolangType(property.type, property.isArray, property.required, property.arrayDimensions, False)}  `yaml:"${property.name}${printOmitemptyIfNeeded(property)}"`
        % endfor
}

        % if needsBuilder(type):
func New${type.name}() *${type.name} {
    return &${type.name}{
            % for property in type.properties:
                % if property.isArray or isinstance(property.type, model.DictionaryType):
        ${stringUtils.toUpperCamelCase(property.name)}: make(${printGolangType(property.type, property.isArray, property.required, property.arrayDimensions, False)}, 0),
                % elif isinstance(property.type, model.ComplexType) and needsBuilder(property.type):
        ${stringUtils.toUpperCamelCase(property.name)}: ${printStarIfRequired(property)}New${property.type.name}(),
                % endif
            % endfor
    }
}
        % endif

func Copy${type.name}(src *${type.name}) *${type.name} {
    if src == nil {
        return nil
    }
    var ret ${type.name}
    % for property in type.properties:
        % if property.isArray:
    ret.${stringUtils.toUpperCamelCase(property.name)} = make(${printGolangType(property.type, property.isArray, property.required, property.arrayDimensions, False)}, 0)
    for _, e := range src.${stringUtils.toUpperCamelCase(property.name)} {
        ret.${stringUtils.toUpperCamelCase(property.name)} = append(ret.${stringUtils.toUpperCamelCase(property.name)}, e)
    }
        % elif isinstance(property.type, model.DictionaryType):
    ret.${stringUtils.toUpperCamelCase(property.name)} = make(${printGolangType(property.type, property.isArray, property.required, property.arrayDimensions, False)}, 0)
    for k, v := range src.${stringUtils.toUpperCamelCase(property.name)} {
        ret.${stringUtils.toUpperCamelCase(property.name)}[k] = v
    }
        % elif isinstance(property.type, model.ComplexType):
    ret.${stringUtils.toUpperCamelCase(property.name)} = ${printAsteriskIfRequired(property)}${copyFunctionName(property)}(${printAndIfRequired(property)}src.${stringUtils.toUpperCamelCase(property.name)})
        % else:
    ret.${stringUtils.toUpperCamelCase(property.name)} = src.${stringUtils.toUpperCamelCase(property.name)}
        % endif
    % endfor

    return &ret
}


    % endif


% endfor
