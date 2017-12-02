package builder_test

import (
	"errors"
	"fmt"

	"github.com/pivotal-cf/kiln/builder"
	"github.com/pivotal-cf/kiln/builder/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MetadataBuilder", func() {
	var (
		iconEncoder                   *fakes.IconEncoder
		logger                        *fakes.Logger
		metadataReader                *fakes.MetadataReader
		releaseManifestReader         *fakes.ReleaseManifestReader
		runtimeConfigsDirectoryReader *fakes.MetadataPartsDirectoryReader
		stemcellManifestReader        *fakes.StemcellManifestReader
		variablesDirectoryReader      *fakes.MetadataPartsDirectoryReader
		formDirectoryReader           *fakes.MetadataPartsDirectoryReader
		instanceGroupDirectoryReader  *fakes.MetadataPartsDirectoryReader
		jobsDirectoryReader           *fakes.MetadataPartsDirectoryReader

		tileBuilder builder.MetadataBuilder
	)

	BeforeEach(func() {
		iconEncoder = &fakes.IconEncoder{}
		logger = &fakes.Logger{}
		metadataReader = &fakes.MetadataReader{}
		releaseManifestReader = &fakes.ReleaseManifestReader{}
		runtimeConfigsDirectoryReader = &fakes.MetadataPartsDirectoryReader{}
		formDirectoryReader = &fakes.MetadataPartsDirectoryReader{}
		instanceGroupDirectoryReader = &fakes.MetadataPartsDirectoryReader{}
		jobsDirectoryReader = &fakes.MetadataPartsDirectoryReader{}
		stemcellManifestReader = &fakes.StemcellManifestReader{}
		variablesDirectoryReader = &fakes.MetadataPartsDirectoryReader{}

		iconEncoder.EncodeReturns("base64-encoded-icon-path", nil)

		releaseManifestReader.ReadStub = func(path string) (builder.ReleaseManifest, error) {
			switch path {
			case "/path/to/release-1.tgz":
				return builder.ReleaseManifest{
					Name:    "release-1",
					Version: "version-1",
				}, nil
			case "/path/to/release-2.tgz":
				return builder.ReleaseManifest{
					Name:    "release-2",
					Version: "version-2",
				}, nil
			default:
				return builder.ReleaseManifest{}, fmt.Errorf("could not read release %q", path)
			}
		}
		runtimeConfigsDirectoryReader.ReadStub = func(path string) ([]builder.Part, error) {
			switch path {
			case "/path/to/runtime-configs/directory":
				return []builder.Part{
					{
						File: "runtime-config-1.yml",
						Name: "runtime-config-1",
						Metadata: map[interface{}]interface{}{
							"name":           "runtime-config-1",
							"runtime_config": "runtime-config-1-manifest",
						},
					},
					{
						File: "runtime-config-2.yml",
						Name: "runtime-config-2",
						Metadata: map[interface{}]interface{}{
							"name":           "runtime-config-2",
							"runtime_config": "runtime-config-2-manifest",
						},
					},
				}, nil
			case "/path/to/other/runtime-configs/directory":
				return []builder.Part{
					{
						File: "runtime-config-3.yml",
						Name: "runtime-config-3",
						Metadata: map[interface{}]interface{}{
							"name":           "runtime-config-3",
							"runtime_config": "runtime-config-3-manifest",
						},
					},
				}, nil
			default:
				return []builder.Part{}, fmt.Errorf("could not read runtime configs directory %q", path)
			}
		}
		formDirectoryReader.ReadStub = func(path string) ([]builder.Part, error) {
			switch path {
			case "/path/to/forms/directory":
				return []builder.Part{
					{
						File: "form-1.yml",
						Name: "form-1",
						Metadata: map[interface{}]interface{}{
							"some-key-1": "some-value-1",
						},
					},
					{
						File: "form-2.yml",
						Name: "form-2",
						Metadata: map[interface{}]interface{}{
							"some-key-2": "some-value-2",
						},
					},
				}, nil
			default:
				return []builder.Part{}, fmt.Errorf("could not read forms directory %q", path)
			}
		}

		instanceGroupDirectoryReader.ReadStub = func(path string) ([]builder.Part, error) {
			switch path {
			case "/path/to/instance-groups/directory":
				return []builder.Part{
					{
						File: "some-instance-group-1.yml",
						Name: "some-instance-group-1",
						Metadata: map[interface{}]interface{}{
							"name": "some-instance-group-1",
							"templates": []interface{}{
								"some-job-1.yml",
							},
						},
					},
					{
						File: "some-instance-group-2.yml",
						Name: "some-instance-group-2",
						Metadata: map[interface{}]interface{}{
							"name": "some-instance-group-2",
							"templates": []interface{}{
								"some-job-2.yml",
							},
						},
					},
				}, nil
			default:
				return []builder.Part{}, fmt.Errorf("could not read instance groups directory %q", path)
			}
		}

		jobsDirectoryReader.ReadStub = func(path string) ([]builder.Part, error) {
			switch path {
			case "/path/to/jobs/directory":
				return []builder.Part{
						{
							File: "some-job-1.yml",
							Name: "some-job-1",
							Metadata: map[interface{}]interface{}{
								"name":    "some-job-1",
								"release": "some-release-1",
							},
						},
						{
							File: "some-job-2.yml",
							Name: "some-job-2",
							Metadata: map[interface{}]interface{}{
								"name":    "some-job-2",
								"release": "some-release-2",
							},
						},
					},
					nil
			default:
				return []builder.Part{}, fmt.Errorf("could not read instance groups directory %q", path)
			}
		}

		variablesDirectoryReader.ReadStub = func(path string) ([]builder.Part, error) {
			switch path {
			case "/path/to/variables/directory":
				return []builder.Part{
					{
						File: "variable-1.yml",
						Name: "variable-1",
						Metadata: map[interface{}]interface{}{
							"name": "variable-1",
							"type": "certificate",
						},
					},
					{
						File: "variable-2.yml",
						Name: "variable-2",
						Metadata: map[interface{}]interface{}{
							"name": "variable-2",
							"type": "user",
						},
					},
				}, nil
			case "/path/to/other/variables/directory":
				return []builder.Part{
					{
						File: "variable-3.yml",
						Name: "variable-3",
						Metadata: map[interface{}]interface{}{
							"name": "variable-3",
							"type": "password",
						},
					},
				}, nil
			default:
				return []builder.Part{}, fmt.Errorf("could not read variables directory %q", path)
			}
		}
		stemcellManifestReader.ReadReturns(builder.StemcellManifest{
			Version:         "2332",
			OperatingSystem: "ubuntu-trusty",
		},
			nil,
		)

		tileBuilder = builder.NewMetadataBuilder(
			formDirectoryReader,
			instanceGroupDirectoryReader,
			jobsDirectoryReader,
			releaseManifestReader,
			runtimeConfigsDirectoryReader,
			variablesDirectoryReader,
			stemcellManifestReader,
			metadataReader,
			logger,
			iconEncoder,
		)
	})

	Describe("Build", func() {
		BeforeEach(func() {
			metadataReader.ReadReturns(builder.Metadata{
				"name":                      "cool-product",
				"metadata_version":          "some-metadata-version",
				"provides_product_versions": "some-provides-product-versions",
				"icon_image":                "unused-icon-image-IGNORE-ME",
				"form_types":                "unused-form-types-IGNORE-ME",
				"job_types":                 "unused-job-types-IGNORE-ME",
			},
				nil,
			)
		})

		It("creates a GeneratedMetadata with the correct information", func() {
			generatedMetadata, err := tileBuilder.Build(builder.BuildInput{
				MetadataPath:             "/some/path/metadata.yml",
				ReleaseTarballs:          []string{"/path/to/release-1.tgz", "/path/to/release-2.tgz"},
				StemcellTarball:          "/path/to/test-stemcell.tgz",
				FormDirectories:          []string{"/path/to/forms/directory"},
				InstanceGroupDirectories: []string{"/path/to/instance-groups/directory"},
				JobDirectories:           []string{"/path/to/jobs/directory"},
				RuntimeConfigDirectories: []string{"/path/to/runtime-configs/directory", "/path/to/other/runtime-configs/directory"},
				VariableDirectories:      []string{"/path/to/variables/directory", "/path/to/other/variables/directory"},
				IconPath:                 "some-icon-path",
				Version:                  "1.2.3",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(stemcellManifestReader.ReadArgsForCall(0)).To(Equal("/path/to/test-stemcell.tgz"))
			metadataPath, version := metadataReader.ReadArgsForCall(0)
			Expect(metadataPath).To(Equal("/some/path/metadata.yml"))
			Expect(version).To(Equal("1.2.3"))

			Expect(generatedMetadata.Name).To(Equal("cool-product"))
			Expect(generatedMetadata.FormTypes).To(Equal([]builder.Part{
				{
					File: "form-1.yml",
					Name: "form-1",
					Metadata: map[interface{}]interface{}{
						"some-key-1": "some-value-1",
					},
				},
				{
					File: "form-2.yml",
					Name: "form-2",
					Metadata: map[interface{}]interface{}{
						"some-key-2": "some-value-2",
					},
				},
			}))
			Expect(generatedMetadata.JobTypes).To(Equal([]builder.Part{
				{
					File: "some-instance-group-1.yml",
					Name: "some-instance-group-1",
					Metadata: map[interface{}]interface{}{
						"name": "some-instance-group-1",
						"templates": []interface{}{
							map[interface{}]interface{}{
								"name":    "some-job-1",
								"release": "some-release-1",
							},
						},
					},
				},
				{
					File: "some-instance-group-2.yml",
					Name: "some-instance-group-2",
					Metadata: map[interface{}]interface{}{
						"name": "some-instance-group-2",
						"templates": []interface{}{
							map[interface{}]interface{}{
								"name":    "some-job-2",
								"release": "some-release-2",
							},
						},
					},
				},
			}))
			Expect(generatedMetadata.Releases).To(Equal([]builder.Release{
				{
					Name:    "release-1",
					Version: "version-1",
					File:    "release-1.tgz",
				},
				{
					Name:    "release-2",
					Version: "version-2",
					File:    "release-2.tgz",
				},
			}))
			Expect(generatedMetadata.RuntimeConfigs).To(Equal([]builder.Part{
				{
					File: "runtime-config-1.yml",
					Name: "runtime-config-1",
					Metadata: map[interface{}]interface{}{
						"name":           "runtime-config-1",
						"runtime_config": "runtime-config-1-manifest",
					},
				},
				{
					File: "runtime-config-2.yml",
					Name: "runtime-config-2",
					Metadata: map[interface{}]interface{}{
						"name":           "runtime-config-2",
						"runtime_config": "runtime-config-2-manifest",
					},
				},
				{
					File: "runtime-config-3.yml",
					Name: "runtime-config-3",
					Metadata: map[interface{}]interface{}{
						"name":           "runtime-config-3",
						"runtime_config": "runtime-config-3-manifest",
					},
				},
			}))
			Expect(generatedMetadata.Variables).To(Equal([]builder.Part{
				{
					File: "variable-1.yml",
					Name: "variable-1",
					Metadata: map[interface{}]interface{}{
						"name": "variable-1",
						"type": "certificate",
					},
				},
				{
					File: "variable-2.yml",
					Name: "variable-2",
					Metadata: map[interface{}]interface{}{
						"name": "variable-2",
						"type": "user",
					},
				},
				{
					File: "variable-3.yml",
					Name: "variable-3",
					Metadata: map[interface{}]interface{}{
						"name": "variable-3",
						"type": "password",
					},
				},
			}))
			Expect(generatedMetadata.StemcellCriteria).To(Equal(builder.StemcellCriteria{
				Version:     "2332",
				OS:          "ubuntu-trusty",
				RequiresCPI: false,
			}))
			Expect(generatedMetadata.Metadata).To(Equal(builder.Metadata{
				"metadata_version":          "some-metadata-version",
				"provides_product_versions": "some-provides-product-versions",
			}))

			Expect(logger.PrintfCall.Receives.LogLines).To(Equal([]string{
				"Read manifest for release release-1",
				"Read manifest for release release-2",
				"Read runtime configs from /path/to/runtime-configs/directory",
				"Read runtime configs from /path/to/other/runtime-configs/directory",
				"Read variables from /path/to/variables/directory",
				"Read variables from /path/to/other/variables/directory",
				"Read manifest for stemcell version 2332",
				"Read forms from /path/to/forms/directory",
				"Read instance groups from /path/to/instance-groups/directory",
				"Read jobs from /path/to/jobs/directory",
				"Read metadata",
			}))

			Expect(iconEncoder.EncodeCallCount()).To(Equal(1))
			Expect(iconEncoder.EncodeArgsForCall(0)).To(Equal("some-icon-path"))

			Expect(generatedMetadata.IconImage).To(Equal("base64-encoded-icon-path"))
		})

		Context("failure cases", func() {
			Context("when the release tarball cannot be read", func() {
				It("returns an error", func() {
					releaseManifestReader.ReadReturns(builder.ReleaseManifest{}, errors.New("failed to read release tarball"))

					_, err := tileBuilder.Build(builder.BuildInput{
						ReleaseTarballs: []string{"release-1.tgz"},
					})
					Expect(err).To(MatchError("failed to read release tarball"))
				})
			})

			Context("when the form directory cannot be read", func() {
				It("returns an error", func() {
					formDirectoryReader.ReadReturns([]builder.Part{}, errors.New("some form error"))

					_, err := tileBuilder.Build(builder.BuildInput{
						FormDirectories: []string{"/path/to/missing/form"},
					})
					Expect(err).To(MatchError(`error reading from form directory "/path/to/missing/form": some form error`))
				})
			})

			Context("when the instance group directory cannot be read", func() {
				It("returns an error", func() {
					instanceGroupDirectoryReader.ReadReturns([]builder.Part{}, errors.New("some instance group error"))

					_, err := tileBuilder.Build(builder.BuildInput{
						InstanceGroupDirectories: []string{"/path/to/missing/instance-groups"},
					})
					Expect(err).To(MatchError(`error reading from instance group directory "/path/to/missing/instance-groups": some instance group error`))
				})
			})

			Context("when the job directory cannot be read", func() {
				It("returns an error", func() {
					jobsDirectoryReader.ReadReturns([]builder.Part{}, errors.New("some job error"))

					_, err := tileBuilder.Build(builder.BuildInput{
						JobDirectories: []string{"/path/to/missing/jobs"},
					})
					Expect(err).To(MatchError(`error reading from job directory "/path/to/missing/jobs": some job error`))
				})
			})

			Context("when an instance group references a job that cannot be found", func() {
				BeforeEach(func() {
					instanceGroupDirectoryReader.ReadStub = func(path string) ([]builder.Part, error) {
						switch path {
						case "/path/to/instance-groups/directory":
							return []builder.Part{
								{
									File: "some-instance-group-1.yml",
									Name: "some-instance-group-1",
									Metadata: map[interface{}]interface{}{
										"name": "some-instance-group-1",
										"templates": []interface{}{
											"some-missing-job",
										},
									},
								},
								{
									File: "some-instance-group-2.yml",
									Name: "some-instance-group-2",
									Metadata: map[interface{}]interface{}{
										"name": "some-instance-group-2",
										"templates": []interface{}{
											"some-job-2",
										},
									},
								},
							}, nil
						default:
							return []builder.Part{}, fmt.Errorf("could not read instance groups directory %q", path)
						}
					}
				})
				It("returns an error", func() {
					_, err := tileBuilder.Build(builder.BuildInput{
						InstanceGroupDirectories: []string{"/path/to/instance-groups/directory"},
						JobDirectories:           []string{"/path/to/jobs/directory"},
					})
					Expect(err).To(MatchError(`instance group "some-instance-group-1" references non-existent job "some-missing-job"`))
				})
			})

			Context("when the runtime configs directory cannot be read", func() {
				It("returns an error", func() {
					runtimeConfigsDirectoryReader.ReadReturns([]builder.Part{}, errors.New("some error"))

					_, err := tileBuilder.Build(builder.BuildInput{
						RuntimeConfigDirectories: []string{"/path/to/missing/runtime-configs"},
					})
					Expect(err).To(MatchError(`error reading from runtime configs directory "/path/to/missing/runtime-configs": some error`))
				})
			})

			Context("when the variables directory cannot be read", func() {
				It("returns an error", func() {
					variablesDirectoryReader.ReadReturns([]builder.Part{}, errors.New("some error"))

					_, err := tileBuilder.Build(builder.BuildInput{
						VariableDirectories: []string{"/path/to/missing/variables"},
					})
					Expect(err).To(MatchError(`error reading from variables directory "/path/to/missing/variables": some error`))
				})
			})

			Context("when the stemcell tarball cannot be read", func() {
				It("returns an error", func() {
					stemcellManifestReader.ReadReturns(builder.StemcellManifest{}, errors.New("failed to read stemcell tarball"))

					_, err := tileBuilder.Build(builder.BuildInput{
						StemcellTarball: "stemcell.tgz",
					})
					Expect(err).To(MatchError("failed to read stemcell tarball"))
				})
			})

			Context("when the icon cannot be encoded", func() {
				BeforeEach(func() {
					iconEncoder.EncodeReturns("", errors.New("failed to encode poncho"))
				})

				It("returns an error", func() {
					_, err := tileBuilder.Build(builder.BuildInput{
						IconPath: "some-icon-path",
					})
					Expect(err).To(MatchError("failed to encode poncho"))
				})
			})

			Context("when the metadata cannot be read", func() {
				It("returns an error", func() {
					metadataReader.ReadReturns(builder.Metadata{}, errors.New("failed to read metadata"))

					_, err := tileBuilder.Build(builder.BuildInput{
						MetadataPath: "metadata.yml",
					})
					Expect(err).To(MatchError("failed to read metadata"))
				})
			})

			Context("when the metadata does not contain a product name", func() {
				It("returns an error", func() {
					metadataReader.ReadReturns(builder.Metadata{
						"metadata_version":          "some-metadata-version",
						"provides_product_versions": "some-provides-product-versions",
					},
						nil,
					)

					_, err := tileBuilder.Build(builder.BuildInput{
						MetadataPath: "metadata.yml",
					})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring(`missing "name" in tile metadata`))
				})
			})

			Context("when the base metadata contains a runtime_configs section", func() {
				It("returns an error", func() {
					metadataReader.ReadReturns(builder.Metadata{
						"name":            "cool-product",
						"runtime_configs": "some-runtime-configs",
					},
						nil,
					)

					_, err := tileBuilder.Build(builder.BuildInput{
						MetadataPath: "metadata.yml",
					})
					Expect(err).To(MatchError("runtime_config section must be defined using --runtime-configs-directory flag"))
				})
			})

			Context("when the base metadata contains a variables section", func() {
				It("returns an error", func() {
					metadataReader.ReadReturns(builder.Metadata{
						"name":      "cool-product",
						"variables": "some-variables",
					},
						nil,
					)

					_, err := tileBuilder.Build(builder.BuildInput{
						MetadataPath: "metadata.yml",
					})
					Expect(err).To(MatchError("variables section must be defined using --variables-directory flag"))
				})
			})
		})
	})
})
